class BusinessView extends KDView

  constructor: ->

    super

    {router} = KD.singletons

    @pricingButton = new KDButtonView
      title       : "See Pricing"
      style       : "solid thin medium thin-white"
      callback    : -> router.handleRoute "/Pricing/Team"

    @signUpButton = new KDButtonView
      title       : "Sign Up Now"
      style       : "solid medium green"
      callback    : -> router.handleRoute "/Register"

    @footer = new FooterView

  viewAppended: JView::viewAppended

  pistachio : ->
    """
      <section class="introduction">
        <div class="inner-container clearfix">
          <article>
            <h2>Koding for Busy People</h2>
            <p>
              Have your private Koding in the cloud, with your rules, your apps and your team.
            </p>
            {{> @signUpButton}}
            {{> @pricingButton}}
          </article>
        </div>
      </section>

      <section class="screenshots">
        <div class="inner-container">
          <figure class="first">
            <img src="/a/images/ss-activity.jpg" alt="Activity">
          </figure>
          <figure class="second">
            <img src="/a/images/ss-terminal.jpg" alt="Terminal">
          </figure>
          <figure class="third">
            <img src="/a/images/ss-environments.jpg" alt="Environments">
          </figure>
        </div>
      </section>

      <section class="features">
        <div class="inner-container clearfix">
          <article class="feature">
            <i class="gameplan icon"></i>
            <h5>Total control over the big picture</h5>
            <p>
              Never miss a thing. Who is working on what, who needs help,
              what needs to be done. Look back into the progress of your
              team’s progress.
            </p>
          </article>
          <article class="feature">
            <i class="ruler icon"></i>
            <h5>Scale as you grow</h5>
            <p>
              Fully scalable environments with customizable stacks gives you
              the ability to scale as you get bigger in size. Size matters.
            </p>
          </article>
          <article class="feature">
            <i class="box-open icon"></i>
            <h5>Ready to roll VM’s</h5>
            <p>
              Stop wasting time setting up environments for every single
              team member as they join in. With a single click,
              they are ready to go.
            </p>
          </article>
          <article class="feature">
            <i class="starflag icon"></i>
            <h5>White Label</h5>
            <p>
              To suit your brand guidelines, fully customisable Koding
              experience in your intranet.
            </p>
          </article>
        </div>
      </section>

      <section class="testimonials">
        <div class="inner-container clearfix">
          <h3 class="general-title">What did they say</h3>
          <h4 class="general-subtitle">People love Koding for a reason. Guess what that reason is?</h4>

          <article>
            <p>It just f***in works! And therefore I love it like I ove my mom.</p>
            <span class="name">JASON FRIEDMANN</span>
          </article>

          <article>
            <p>It just f***in works! And therefore I love it like I ove my mom.</p>
            <span class="name">JASON FRIEDMANN</span>
          </article>

          <article>
            <p>It just f***in works! And therefore I love it like I ove my mom.</p>
            <span class="name">JASON FRIEDMANN</span>
          </article>

          <article>
            <p>It just f***in works! And therefore I love it like I ove my mom.</p>
            <span class="name">JASON FRIEDMANN</span>
          </article>
        </div>
      </section>

      <section class='check-out'>
        <h3><a href='/Pricing'>Check out our pricing</a> and get started with Koding right away!</h3>
      </section>
      {{> @footer}}
    """


